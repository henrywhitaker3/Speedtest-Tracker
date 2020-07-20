<?php

namespace App;

use App\Helpers\SettingsHelper;
use Illuminate\Database\Eloquent\Model;

class Setting extends Model
{
    /**
     * The attributes that are mass assignable.
     *
     * @var array
     */
    protected $fillable = [
        'name', 'value', 'description'
    ];

    protected $table = 'settings';
}
