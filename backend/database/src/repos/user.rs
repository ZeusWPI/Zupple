use sqlx::SqlitePool;

use crate::{
    error::DatabaseError,
    models::user::{User, UserCreate, UserPatch},
};

pub struct UserRepo<'a> {
    db: &'a SqlitePool,
}

impl<'a> UserRepo<'a> {
    pub fn new(db: &'a SqlitePool) -> Self {
        Self { db }
    }

    pub async fn by_id(&self, id: u32) -> Result<User, DatabaseError> {
        sqlx::query_as("SELECT id, username FROM user WHERE id == ? LIMIT 1;")
            .bind(id)
            .fetch_optional(self.db)
            .await?
            .ok_or(DatabaseError::NotFound)
    }

    pub async fn create(&self, user: UserCreate) -> Result<User, DatabaseError> {
        sqlx::query_as(
            "
            INSERT INTO user (username) VALUES (?)
            RETURNING id, username;
            ",
        )
        .bind(&user.username)
        .bind(&user.username)
        .fetch_optional(self.db)
        .await?
        .ok_or(DatabaseError::NotFound)
    }

    pub async fn patch(&self, user_id: u32, patch_user: UserPatch) -> Result<User, DatabaseError> {
        sqlx::query_as(
            "
            UPDATE user SET username = ? WHERE id = ?
            RETURNING id, username;
            ",
        )
        .bind(patch_user.username)
        .bind(user_id)
        .fetch_optional(self.db)
        .await?
        .ok_or(DatabaseError::NotFound)
    }
}
