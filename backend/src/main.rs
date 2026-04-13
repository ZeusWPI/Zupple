use std::error::Error;
use tracing_subscriber::{EnvFilter, fmt, prelude::*};
use zupple::{config::AppConfig, puzzle::kuromasu::Kuromasu, start_app};

#[tokio::main]
async fn main() -> Result<(), Box<dyn Error>> {
    let _ = dotenvy::dotenv();
    let config = AppConfig::load()?;


    // setup logging
    tracing_subscriber::registry()
        .with(fmt::layer())
        .with(EnvFilter::from_env("LOG_LEVEL"))
        .init();

    let kuro = Kuromasu::new(8)?;
    println!("{}", kuro);

    start_app(config).await?;

    Ok(())
}
