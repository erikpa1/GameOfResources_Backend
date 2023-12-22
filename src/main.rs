#![allow(warnings)]

use std::env;
use std::thread;

use std::sync::Mutex;
use tokio::signal;

use actix_files as fs;
use actix_web::middleware::Logger;
use actix_web::web::ServiceConfig;
use actix_web::{get, post, web, App, HttpResponse, HttpServer, Responder};

//Toto musi byt kvoli start v App
use actix::prelude::*;

//Bez tohto actix funkcie nevypisuju
use env_logger;

use async_std::task;

mod app;
mod warehouse;

#[get("/api/test")]
async fn _test_function(data: web::Data<Mutex<app::App>>) -> impl Responder {
    let tmp = data.lock().unwrap();

    let resCount = tmp.warehouse.get_resource_count("gold");

    HttpResponse::Ok().body(format!("{}", resCount))
    // HttpResponse::Ok().body("Ok")
}

pub fn mount(reference: &mut ServiceConfig) {
    reference.service(_test_function);
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    std::env::set_var("RUST_LOG", "actix_web=info");

    env_logger::init();

    let mut app = app::App::new();
    // let mut app = app.start();

    let mut appData = web::Data::new(Mutex::new(app));

    if true == true {
        thread::spawn(|| {
            println!("Spawned tokio thread for Ctrl+C interuption");
            task::block_on(async {
                loop {
                    if let Ok(_) = signal::ctrl_c().await {
                        println!("End signal Hit");
                    }
                }
            });

            // Print whenever a HUP signal is received
        });
    }

    HttpServer::new(move || {
        App::new()
            .wrap(Logger::default())
            .app_data(web::Data::clone(&appData))
            .app_data(web::FormConfig::default())
            .configure(mount)
            .service(
                fs::Files::new("/", "static")
                    .index_file("index.html")
                    .show_files_listing(),
            )
    })
    .bind(("0.0.0.0", 5000))?
    .run()
    .await
}