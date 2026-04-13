use sqlx::{SqlitePool, migrate::MigrateDatabase, sqlite::SqlitePoolOptions};

use crate::{error::DatabaseError, repos::user::UserRepo};

pub mod models {
    pub mod user;
}

pub mod repos {
    pub mod user;
}

pub mod error;

#[derive(Clone)]
pub struct Database {
    db: SqlitePool,
}

impl Database {
    pub fn new(db: SqlitePool) -> Self {
        Self { db }
    }

    pub async fn create_connect_migrate(db_url: &str) -> Result<Self, DatabaseError> {
        // create database if not exists
        sqlx::Sqlite::create_database(db_url).await?;

        // connect to database
        let db = SqlitePoolOptions::new()
            .max_connections(5)
            .connect(db_url)
            .await?;

        // run migrations
        sqlx::migrate!("../migrations").run(&db).await?;

        Ok(Self { db })
    }

    pub fn users<'a>(&'a self) -> UserRepo<'a> {
        UserRepo::new(&self.db)
    }
}
