<?php

use CodeIgniter\Router\RouteCollection;

/**
 * @var RouteCollection $routes
 */
$routes->get('/', 'Home::index');

$routes->post("/notifications", "Producer::SendMessage");
// $routes->resource('Producer');


