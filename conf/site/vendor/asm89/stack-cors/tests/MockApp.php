<?php

namespace Asm89\Stack\Tests;

use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\HttpKernel\HttpKernelInterface;

class MockApp implements HttpKernelInterface
{

    private $responseHeaders;

    public function __construct(array $responseHeaders)
    {
        $this->responseHeaders = $responseHeaders;
    }

    public function handle(Request $request, $type = HttpKernelInterface::MASTER_REQUEST, $catch = true)
    {
        $response = new Response();

        $response->headers->add($this->responseHeaders);

        return $response;
    }
}
