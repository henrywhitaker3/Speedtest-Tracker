<?php

namespace App\Models\Auth;

use Illuminate\Database\Eloquent\Model;

class LoginSession extends Model
{
    protected $fillable = [
        'id',
        'token',
        'active',
        'user_id',
        'expires',
        'ip'
    ];

    protected $table = 'active_sessions';

    public $incrementing = false;
}
