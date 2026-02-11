Race Condition Simulation — Go + PostgreSQL

📌 Overview

Repository ini dibuat untuk mensimulasikan race condition pada pola:

    Update child
    Recalculate aggregate
    Update parent

Kasus ini merepresentasikan situasi di mana dua proses paralel:

- Mengupdate child row

- Melakukan agregasi

- Mengupdate parent row

Tanpa mekanisme locking, race condition dapat menyebabkan parent tidak pernah ter-update meskipun semua child sudah memenuhi kondisi.

---

🧠 Problem Statement

Ketika dua transaksi berjalan hampir bersamaan:

1. Transaksi A mengupdate child 1 → APPROVED

2. Transaksi B mengupdate child 2 → APPROVED

3. Keduanya membaca snapshot data yang belum lengkap

4. Keduanya menyimpulkan belum semua approved

5. Parent tidak pernah menjadi APPROVED

Masalah ini terjadi karena:

- Isolation level default PostgreSQL: READ COMMITTED

- Tidak adanya SELECT ... FOR UPDATE

- Tidak ada koordinasi antar transaksi

---

✅ What This Repo Demonstrates

- ❌ Race condition tanpa locking

- ✅ Konsistensi dengan SELECT ... FOR UPDATE

- ✅ Perilaku blocking row-level lock PostgreSQL

- ✅ Concurrency simulation dengan goroutine

---

🛠 Requirements

- Go 1.18+

- PostgreSQL (local atau via SSH tunnel)

- psql (optional untuk testing manual)

---

🗄 Database Setup

1️⃣ Create Schema and Tables

    CREATE SCHEMA IF NOT EXISTS demo;

    DROP TABLE IF EXISTS demo.child;
    DROP TABLE IF EXISTS demo.history;

    CREATE TABLE demo.history (
        id VARCHAR PRIMARY KEY,
        result VARCHAR
    );

    CREATE TABLE demo.child (
        id VARCHAR PRIMARY KEY,
        history_id VARCHAR NOT NULL,
        result VARCHAR,
        CONSTRAINT fk_history
            FOREIGN KEY(history_id)
            REFERENCES demo.history(id)
    );

    CREATE INDEX idx_history_id
    ON demo.child(history_id);

---

2️⃣ Seed Data

    INSERT INTO demo.history (id, result)
    VALUES ('H1', NULL);

    INSERT INTO demo.child (id, history_id, result)
    VALUES
    ('C1', 'H1', NULL),
    ('C2', 'H1', NULL);

---

🚀 Setup & Run

1️⃣ Initialize Go Module

    go mod init race-demo

2️⃣ Install Dependency

    go get github.com/lib/pq
    go mod tidy

---

⚙ Configure Database Connection

Edit DSN di main.go:

    const dsn = "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"

Jika menggunakan SSH tunnel:

    ssh -L 5433:localhost:5432 user@server

Maka ubah DSN:

    port=5432

---

🧪 Running the Simulation

🔴 Without Lock (Race Condition Mode)

    go run main.go ./without-lock

Kemungkinan hasil:

    Final History Result:

Parent tetap NULL meskipun semua child sudah APPROVED.

Race condition terjadi ✅

---

🟢 With Lock

Jika menggunakan mode lock:

    go run main.go ./with-lock

Output:

    History updated by S2
    Final History Result: APPROVED

Race condition hilang ✅

---

🔍 How It Works

Tanpa Lock:

- Dua transaksi berjalan paralel

- Masing-masing tidak melihat update transaksi lain

- Snapshot isolation menyebabkan agregasi salah

Dengan Lock:

    SELECT id FROM demo.history WHERE id='H1' FOR UPDATE;

- Transaksi kedua menunggu

- Eksekusi menjadi serialized per history_id

- Agregasi selalu konsisten

---

📊 What To Experiment

- Tambahkan jumlah goroutine

- Tambahkan delay lebih besar

- Loop 100x untuk melihat probabilistic race

- Monitor lock di Postgres:

  SELECT \* FROM pg_stat_activity WHERE wait_event_type='Lock';

---

🧩 Key Takeaways

- Race condition bukan bug logic biasa

- Masalah terjadi pada coordination boundary

- Parent row adalah boundary yang tepat untuk locking

- Row-level locking PostgreSQL sangat granular dan scalable

- FOR UPDATE ≠ table lock

---

📚 Concepts Covered

- MVCC (Multi-Version Concurrency Control)

- Transaction Isolation Level (READ COMMITTED)

- Row-level Lock

- Critical Section

- Concurrency Simulation in Go

- Aggregation Consistency Problem

---

⚠ Important Notes

- Race condition bersifat timing-dependent

- Kadang tidak muncul tanpa artificial delay

- Jangan gunakan time.Sleep di production (hanya untuk simulasi)

- Pastikan ada index pada foreign key (history_id)

---
