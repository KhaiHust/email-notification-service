# 📧 Email Notification Management System

A centralized email notification platform for internal development teams, designed to unify email delivery across services, improve operational efficiency, and reduce technical debt.
---

## 🚀 Features

- ✅ Manage dynamic email templates (with variables)
- ✅ Add and manage custom SMTP providers
- ✅ Asynchronous email delivery using Kafka
- ✅ Scheduled email sending with Google Cloud Tasks
- ✅ Auto retry and deduplication to ensure idempotency
- ✅ Email open tracking via tracking pixel
- ✅ RESTful API for system integration
---

## 🛠️ Tech Stack

| Layer       | Technologies |
|-------------|--------------|
| **Backend** | Golang (Gin), Redis, PostgreSQL, Kafka |
| **Infra**   | Docker, Google Cloud Tasks, GitOps |
| **Monitoring** | NewRelic, Grafana, Promtail |
| **Test**    | Load test with k6 |

---

## ⚙️ Architecture

### 📐 System Design
- Microservices Architecture
- Hexagonal Architecture (Ports & Adapters)
- Pub/Sub Pattern for decoupling (Kafka)
- Strategy Pattern for SMTP provider integration

### 📦 Core Services
- **Public Service**: UI and workspace management
- **Internal Service**: API integration with other services
- **Worker Service**: Background processing (sending emails)
- **Migration Service**: DB schema updates

---

## 📊 Performance

- ✅ Load tested with 2300+ requests using [k6](https://k6.io)
- ⏱️ P95 latency: 94.53 ms
- 📬 Email success rate: 99.99%
- 🧩 Supports horizontal scaling of worker nodes

---

## 🧠 Key Design Decisions

| Problem | Solution |
|--------|----------|
| Email delay & blocking | Use Kafka for async processing |
| Duplicate delivery (retry/crash) | Redis-based deduplication by `email_request_id` |
| Cronjob inaccuracy | Switched to Google Cloud Tasks for precision and retry |
| Poor observability | Integrated NewRelic APM & centralized logging |

---

## 📌 Future Improvements

- [ ] Add bounce & unsubscribe handling via webhook
- [ ] Test coverage for critical components
- [ ] Add email click tracking

---

## 📄 License

This project is academic work. Not intended for production use without security & compliance checks.

---

## 🙌 Acknowledgments

Thanks to my supervisor, industry mentors, and friends who supported this project.

