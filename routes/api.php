<?php

use App\Http\Controllers\SpeedtestController;
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
    'middleware' => [ 'api' ],
    'prefix' => 'speedtest'
], function($router) {
    Route::get('/', 'SpeedtestController@index')
         ->name('speedtest.index');
    Route::get('latest', 'SpeedtestController@latest')
         ->name('speedtest.latest');
    Route::get('time/{time}', 'SpeedtestController@time')
         ->name('speedtest.time');
    Route::get('run', 'SpeedtestController@run')
         ->name('speedtest.run');
});

Route::group([
    'middleware' => 'api'
], function () {
    Route::get('backup', 'BackupController@backup')
         ->name('data.backup');
    Route::post('restore', 'BackupController@restore')
         ->name('data.restore');
});

Route::group([
    'middleware' => 'api',
    'prefix' => 'update',
], function () {
    Route::get('check', 'UpdateController@checkForUpdate')
         ->name('update.check');
    Route::get('download', 'UpdateController@downloadUpdate')
         ->name('update.download');
    Route::get('extract', 'UpdateController@extractUpdate')
         ->name('update.extract');
    Route::get('move', 'UpdateController@moveUpdate')
         ->name('update.move');
});
