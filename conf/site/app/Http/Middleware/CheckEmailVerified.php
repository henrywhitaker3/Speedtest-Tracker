<?php

namespace App\Http\Middleware;

use App\Helpers\EmailVerificationHelper;
use Closure;
use Exception;

class CheckEmailVerified
{
    /**
     * Handle an incoming request.
     *
     * @param  \Illuminate\Http\Request  $request
     * @param  \Closure  $next
     * @return mixed
     */
    public function handle($request, Closure $next)
    {
        try {
            if(!EmailVerificationHelper::userIsVerified()) {
                return response()->json([
                    'error' => 'You need to verify your email address',
                ], 401);
            }
        } catch(Exception $e) {
            return response()->json([
                'error' => 'Your account was not found'
            ], 422);
        }

        return $next($request);
    }
}
