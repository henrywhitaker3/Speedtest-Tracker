<?php

use App\Helpers\SpeedtestHelper;
use App\Http\Controllers\AuthController;
use App\Http\Controllers\BackupController;
use App\Http\Controllers\HomepageDataController;
use App\Http\Controllers\SettingsController;
use App\Http\Controllers\SpeedtestController;
use App\Http\Controllers\UpdateController;
use App\Models\Speedtest;
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

Route::group([
    'middleware' => ['api'],
    'prefix' => 'speedtest'
], function ($router) {
    Route::get('/', [SpeedtestController::class, 'index'])
        ->name('speedtest.index');
    Route::get('latest', [SpeedtestController::class, 'latest'])
        ->name('speedtest.latest');
    Route::get('time/{time}', [SpeedtestController::class, 'time'])
        ->name('speedtest.time');
    Route::get('fail/{time}', [SpeedtestController::class, 'fail'])
        ->name('speedtest.fail');
    Route::get('run', [SpeedtestController::class, 'run'])
        ->name('speedtest.run');
    Route::get('home/{time}', HomepageDataController::class)
        ->name('speedtest.home');

    Route::group([
        'prefix' => 'delete'
    ], function () {
        Route::delete('all', [SpeedtestController::class, 'deleteAll']);
        Route::delete('{speedtest}', [SpeedtestController::class, 'delete']);
    });
});

Route::group([
    'middleware' => 'api'
], function () {
    Route::get('backup', [BackupController::class, 'backup'])
        ->name('data.backup');
    Route::post('restore', [BackupController::class, 'restore'])
        ->name('data.restore');
});

Route::group([
    'middleware' => 'api',
    'prefix' => 'update',
], function () {
    Route::get('changelog', [UpdateController::class, 'changelog'])
        ->name('update.changelog');
    Route::get('check', [UpdateController::class, 'checkForUpdate'])
        ->name('update.check');
    Route::get('download', [UpdateController::class, 'downloadUpdate'])
        ->name('update.download');
    Route::get('extract', [UpdateController::class, 'extractUpdate'])
        ->name('update.extract');
    Route::get('move', [UpdateController::class, 'moveUpdate'])
        ->name('update.move');
});

Route::group([
    'middleware' => 'api',
    'prefix' => 'settings'
], function () {
    Route::get('/config', [SettingsController::class, 'config'])
        ->name('settings.config');
    Route::get('/test-notification', 'IntegrationsController@testNotification')
        ->name('settings.test_notification');
    Route::get('/test-healthchecks/{method}', 'IntegrationsController@testHealthchecks')
        ->name('settings.test_notification');
    Route::get('/', [SettingsController::class, 'index'])
        ->name('settings.index');
    Route::put('/', [SettingsController::class, 'store'])
        ->name('settings.store');
    Route::post('/', [SettingsController::class, 'store'])
        ->name('settings.update');
    Route::post('/bulk', [SettingsController::class, 'bulkStore'])
        ->name('settings.bulk.update');
});

Route::group(
    [
        'middleware' => 'api',
        'prefix' => 'auth'
    ],
    function ($router) {
        Route::post('register', [AuthController::class, 'register'])->name('auth.register');
        Route::post('login', [AuthController::class, 'login'])->middleware('throttle:60,1')->name('auth.login');
        Route::get('logout', [AuthController::class, 'logout'])->name('auth.logout');
        Route::get('refresh', [AuthController::class, 'refresh'])->middleware('throttle:60,1')->name('auth.refresh');
        Route::get('me', [AuthController::class, 'me'])->middleware('session_active')->name('auth.me');
        Route::post('change-password', [AuthController::class, 'changePassword'])->middleware('session_active')->name('auth.change_password');

        Route::group(
            [
                'middleware' => ['api', 'session_active'],
                'prefix' => 'sessions'
            ],
            function ($router) {
                Route::get('/', [AuthController::class, 'getSessions'])->name('auth.sessions.all');
                Route::delete('/{id}', [AuthController::class, 'deleteSession'])->name('auth.sessions.delete');
            }
        );
    }
);
