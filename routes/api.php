<?php

use Illuminate\Http\Request;
use Illuminate\Support\Facades\Route;

/*
|--------------------------------------------------------------------------
| API Routes
|--------------------------------------------------------------------------
|
| Here is where you can register API routes for your application. These
| routes are loaded by the RouteServiceProvider within a group which
| is assigned the "api" middleware group. Enjoy building your API!
|
*/

Route::group(
    [
        'middleware' => 'api',
        'prefix' => 'auth'
    ],
    function ($router) {
        Route::post('register', 'AuthController@register')->name('auth.register');

        Route::get('verify-email', 'AuthController@verifyEmail')->middleware('throttle:5,1')->name('auth.verify_email');

        Route::post('login', 'AuthController@login')->middleware('throttle:10,1')->name('auth.login');

        Route::get('logout', 'AuthController@logout')->name('auth.logout');

        Route::get('refresh', 'AuthController@refresh')->middleware(['throttle:5,1' ])->name('auth.refresh');

        Route::get('me', 'AuthController@me')->middleware(['session_active' ])->name('auth.me');

        Route::put('details', 'UserController@update')->name('auth.user.update_details');

        Route::group(
            [
                'middleware' => ['api', 'session_active' ],
                'prefix' => 'sessions'
            ],
            function($router) {
                Route::get('/', 'AuthController@getSessions')->name('auth.sessions.all');
                Route::delete('/{id}', 'AuthController@deleteSession')->name('auth.sessions.delete');
            }
        );
    }
);

Route::group([
    'middleware' => [ 'api', 'session_active' ],
    'prefix' => 'speedtest'
], function($router) {
    Route::get('latest', 'SpeedtestController@latest')->name('speedtest.latest');
    Route::get('run', 'SpeedtestController@run')->name('speedtest.run');
});
