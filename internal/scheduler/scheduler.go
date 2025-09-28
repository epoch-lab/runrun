package scheduler

import (
	"log"
	"math/rand"
	"runrun/internal"
	"runrun/internal/protocol"
	"time"
)

// Scheduler 定时任务调度器
type Scheduler struct {
	stopChan chan bool
	running  bool
}

// NewScheduler 创建新的调度器实例
func NewScheduler() *Scheduler {
	return &Scheduler{
		stopChan: make(chan bool),
		running:  false,
	}
}

// Start 启动定时任务调度器
func (s *Scheduler) Start() {
	if s.running {
		log.Println("Scheduler is already running")
		return
	}

	s.running = true
	log.Println("Starting scheduler...")

	go s.run()
}

// Stop 停止定时任务调度器
func (s *Scheduler) Stop() {
	if !s.running {
		return
	}

	log.Println("Stopping scheduler...")
	s.stopChan <- true
	s.running = false
}

// run 调度器主循环
func (s *Scheduler) run() {
	// 计算下次执行时间（明天凌晨2点）
	nextRun := getNextRunTime()
	log.Printf("Next daily check scheduled at: %s", nextRun.Format("2006-01-02 15:04:05"))

	timer := time.NewTimer(time.Until(nextRun))
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			// 执行每日检查任务
			s.executeDailyCheck()

			// 设置下次执行时间（明天凌晨2点）
			nextRun = getNextRunTime()
			log.Printf("Next daily check scheduled at: %s", nextRun.Format("2006-01-02 15:04:05"))
			timer.Reset(time.Until(nextRun))

		case <-s.stopChan:
			log.Println("Scheduler stopped")
			return
		}
	}
}

// getNextRunTime 计算下次执行时间（明天凌晨2点）
func getNextRunTime() time.Time {
	now := time.Now()
	// 明天凌晨2点
	next := time.Date(now.Year(), now.Month(), now.Day()+1, 2, 0, 0, 0, now.Location())
	return next
}

// executeDailyCheck 执行每日检查任务
func (s *Scheduler) executeDailyCheck() {
	log.Println("Starting daily check...")

	// 查询所有需要跑步的用户
	var users []internal.User
	result := internal.DB.Where("is_running_required = ?", true).Find(&users)
	if result.Error != nil {
		log.Printf("Failed to query users: %v", result.Error)
		return
	}

	log.Printf("Found %d users requiring running tasks", len(users))

	// 为每个用户安排随机时间的跑步任务
	for _, user := range users {
		s.scheduleUserRun(user)
	}

	log.Println("Daily check completed")
}

// scheduleUserRun 为单个用户安排跑步任务
func (s *Scheduler) scheduleUserRun(user internal.User) {
	// 计算随机延迟时间（10:00-22:00之间，即8-20小时后）
	minHours := 8  // 10:00 - 02:00 = 8小时
	maxHours := 20 // 22:00 - 02:00 = 20小时
	delayHours := minHours + rand.Intn(maxHours-minHours)
	delayMinutes := rand.Intn(60) // 随机分钟数

	delay := time.Duration(delayHours)*time.Hour + time.Duration(delayMinutes)*time.Minute
	runTime := time.Now().Add(delay)

	log.Printf("Scheduled run for user %s at %s (delay: %s)",
		user.Account, runTime.Format("15:04:05"), delay)

	// 使用定时器设置延迟执行
	time.AfterFunc(delay, func() {
		s.runForUser(user)
	})
}

// runForUser 为用户执行一次跑步
func (s *Scheduler) runForUser(user internal.User) {
	log.Printf("Executing run for user: %s", user.Account)

	// 生成随机客户端信息
	clientInfo := protocol.GenerateFakeClient()

	// 登录获取用户信息
	userInfo, err := protocol.Login(user.Account, user.Password, clientInfo)
	if err != nil {
		log.Printf("Failed to login for user %s: %v", user.Account, err)
		return
	}

	// 生成随机跑步参数 (4-5km, 20-30分钟)
	distance := int64(4000 + rand.Intn(1000)) // 4000-4999米
	duration := int32(30 + rand.Intn(10))     // 30-39分钟

	// 提交跑步记录
	err = protocol.Submit(*userInfo, clientInfo, duration, distance)
	if err != nil {
		log.Printf("Failed to submit run for user %s: %v", user.Account, err)
		return
	}

	log.Printf("Successfully submitted run for user %s: %.2fkm in %d minutes",
		user.Account, float64(distance)/1000.0, duration)

	// 更新数据库中的跑步距离
	err = s.updateUserRunningProgress(user.ID, float64(distance)/1000.0)
	if err != nil {
		log.Printf("Failed to update running progress for user %s: %v", user.Account, err)
		return
	}
}

// updateUserRunningProgress 更新用户跑步进度
func (s *Scheduler) updateUserRunningProgress(userID uint, distanceKm float64) error {
	var user internal.User
	result := internal.DB.First(&user, userID)
	if result.Error != nil {
		return result.Error
	}

	// 更新当前距离
	user.CurrentDistance += distanceKm

	// 检查是否达到目标距离
	if user.CurrentDistance >= user.TargetDistance {
		user.IsRunningRequired = false
		log.Printf("User %s has reached target distance %.2f km, stopping automatic runs",
			user.Account, user.TargetDistance)
	}

	// 保存更新
	return internal.DB.Save(&user).Error
}
