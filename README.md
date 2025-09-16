**Helcare Backend (Work-in-Progress)**

Helcare is a proposed healthcare platform designed to connect patients, doctors, and admins. The idea was to create a system where:

Parents can get medical aid, drugs, and nutritional recommendations for their kids.

Doctors can interact with patients, manage appointments, and prescribe treatments.

Admins can oversee doctors, patients, and system-wide activities.

The system was meant to support Web + Mobile clients, with a long-term plan of integrating AI-driven nutritional recommendations for users.

I worked on the backend system (Golang + Fiber) for about 60% of the project. While I didn’t finish all features, I built core modules and learned a lot in the process.

**✅ Features Implemented**

Authentication & Authorization using JWT

Password Reset Flow with secure, token-based mechanism

Implemented an appointment booking system for users to book appointments with doctors

Implemented a cart and billing system for usres to be able to use to purchase drugs off of the platform 

Wallet System with top-up & withdrawal flows integrated with Paystack

Transaction PIN validation

Balance checks & pending balance logic

Webhook handling (success, failure, reversal)

Idempotency & retry-safe design

Database Transactions to ensure financial consistency

Error Logging & Observability groundwork

**🔑 Key Learnings**

Designing clean authentication flows (including secure password resets).

Handling real-money operations safely with atomic DB transactions.

Implementing idempotency & retry logic with third-party APIs (Paystack).

Structuring backend code for reusability and separation of concerns.

The importance of observability (logs & alerts) in financial/healthcare systems.

**⚠️ Status**

This project is currently paused at ~60% completion. The repo is being archived as a portfolio piece.
Instead of finishing, I chose to document what I built and what I learned, and move on to projects that better align with my day-to-day use and long-term goals.

**🛠️ Tech Stack**

Language: Go (Golang)

Framework: Fiber

Database: PostgreSQL

Payments: Paystack API

Auth: JWT-based
