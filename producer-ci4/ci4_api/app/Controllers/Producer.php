<?php
 
namespace App\Controllers;
 
use CodeIgniter\RESTful\ResourceController;
use CodeIgniter\API\ResponseTrait;
use App\Models\ProductModel;
use PhpAmqpLib\Connection\AMQPStreamConnection;

class Producer extends ResourceController
{
    use ResponseTrait;
    public function SendMessage() {
        $data =  $this->request->getVar('message');
        

        $connection = new AMQPStreamConnection('localhost', 5672, 'guest', 'guest');
        $channel = $connection->channel();
        // $message = $this->post('message');

        $channel->exchange_declare('notifications', 'topic', true, false, false);

        $channel->basic_publish(new \PhpAmqpLib\Message\AMQPMessage($data), 'notifications',
            'email.user_registered_on_event');
    
        
        
        $channel->close();
        $connection->close();

        $response = [
            'status' => 201,
            'error' => null,
            'messages' => [
                'success' => 'messages send to rabbitmq',
            ]
            ];
        
        
        return $this->respondCreated($response);
    }
}
