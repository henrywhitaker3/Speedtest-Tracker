<?php

namespace App;

use Illuminate\Database\Eloquent\Model;

class Speedtest extends Model
{
    /**
     * The attributes that are mass assignable.
     *
     * @var array
     */
    protected $fillable = [
        'ping', 'download', 'upload', 'created_at'
    ];

    protected $table = 'speedtests';
}
