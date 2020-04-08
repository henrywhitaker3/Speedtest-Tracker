<?php

namespace App\Auth;

use Illuminate\Database\Eloquent\Model;

class EmailVerification extends Model
{
    protected $fillable = [
        'user_id',
        'token',
        'expires'
    ];
}
