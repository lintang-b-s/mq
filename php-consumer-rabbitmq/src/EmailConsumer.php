<?php
require_once __DIR__ . '/../vendor/autoload.php';

use PhpAmqpLib\Connection\AMQPStreamConnection;

$connection = new AMQPStreamConnection("localhost", 5672, "guest", "guest", "/");

$channel = $connection->channel();


list($queue_name, ,)  = $channel->queue_declare('email-queue', false, true, false, false);

$channel->queue_bind('email-queue', 'notifications', 'email.user_registered_on_event');


echo " [*] menunggu message dikirim oleh producer. \n";

$callback = function($msg) {
    echo ' [x] ', $msg->getRoutingKey(), ':', $msg->getBody(), "\n";
};

$channel->basic_consume('email-queue', '', false, true, false, false, $callback);

try {
    $channel->consume();
} catch (\Throwable $exception) {
    echo $exception->getMessage();
}







