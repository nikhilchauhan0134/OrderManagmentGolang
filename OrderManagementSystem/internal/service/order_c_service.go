package service

import (
	"OrderManagementSystem/internal/concurrency"
	"OrderManagementSystem/internal/contracts"
	"OrderManagementSystem/internal/initresources"
	"OrderManagementSystem/internal/models"
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

type OrderConcurrencyService struct {
	con contracts.OrderCurrencyRespository

	// Semaphore → limit DB concurrency
	dbSemaphore chan struct{}

	// Cache
	cache      map[string]string
	cacheMutex sync.RWMutex

	// Object reuse
	orderPool sync.Pool

	// Inventory control
	inventory     int
	inventoryCond *sync.Cond

	// Metrics
	totalOrders int64

	// ✅ Global Worker Pool
	workerPool *concurrency.WorkerPool
}

// ---------------- CONSTRUCTOR ----------------

func NewOrderConcurrencyService(
	con contracts.OrderCurrencyRespository,
) *OrderConcurrencyService {

	initresources.InitAll()

	mu := &sync.Mutex{}

	wp := concurrency.NewWorkerPool(
		3,                    // min workers
		10,                   // max workers
		100*time.Millisecond, // rate limit
	)

	wp.Start(context.Background())

	return &OrderConcurrencyService{
		con:         con,
		dbSemaphore: make(chan struct{}, 5),
		cache:       make(map[string]string),
		orderPool: sync.Pool{
			New: func() interface{} {
				return new(models.Orders)
			},
		},
		inventoryCond: sync.NewCond(mu),
		inventory:     10,
		workerPool:    wp,
	}
}

// ---------------- WORKER POOL ----------------

// ---------------- MAIN API ----------------

func (s *OrderConcurrencyService) V1CreateOrder(
	ctx context.Context,
	order models.Orders,
) (models.CommonResponse, error) {

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	atomic.AddInt64(&s.totalOrders, 1)

	// sync.Pool
	pooledOrder := s.orderPool.Get().(*models.Orders)
	defer s.orderPool.Put(pooledOrder)
	*pooledOrder = order

	// Semaphore (DB protection)
	select {
	case s.dbSemaphore <- struct{}{}:
		defer func() { <-s.dbSemaphore }()
	case <-ctx.Done():
		return models.CommonResponse{}, ctx.Err()
	}

	// Submit jobs
	s.workerPool.Submit(func(ctx context.Context) error {
		return s.con.CreateOrderConcurrent(ctx, *pooledOrder)
	})

	s.workerPool.Submit(func(ctx context.Context) error {
		return s.reserveInventory(ctx)
	})

	s.workerPool.Submit(func(ctx context.Context) error {
		return s.updateCache(ctx, order.ID)
	})

	// Wait for 3 results
	for i := 0; i < 3; i++ {
		select {
		case err := <-s.workerPool.Results():
			if err != nil {
				return models.CommonResponse{}, err
			}
		case <-ctx.Done():
			return models.CommonResponse{}, ctx.Err()
		}
	}

	return models.CommonResponse{
		Status:  1,
		Message: "Order created successfully",
	}, nil
}

// ---------------- HELPERS ----------------

// Inventory reservation (Mutex + Cond)
func (s *OrderConcurrencyService) reserveInventory(ctx context.Context) error {
	s.inventoryCond.L.Lock()
	defer s.inventoryCond.L.Unlock()

	for s.inventory == 0 {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			s.inventoryCond.Wait()
		}
	}

	s.inventory--
	s.inventoryCond.Signal() // 🔔 wake waiting goroutine
	return nil
}
func (s *OrderConcurrencyService) updateCache(
	ctx context.Context,
	orderID string,
) error {

	if orderID == "" {
		return errors.New("invalid order id")
	}

	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	s.cache[orderID] = "CREATED"
	return nil
}
