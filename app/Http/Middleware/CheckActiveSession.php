<?php

namespace App\Http\Middleware;

use App\Auth\LoginSession;
use Closure;
use Exception;

class CheckActiveSession
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
            $token = $request->bearerToken();
            if($token == null) { $token = $request->token; }
            $session = LoginSession::where('token', $token)
                                   ->first();

            if(!$session->active) {
                return response()->json([
                    'error' => 'token is invalid'
                ], 401);
            }
        } catch(Exception $e) {
            return response()->json([
                'error' => 'token not found'
            ], 401);
        }

        return $next($request);
    }
}
