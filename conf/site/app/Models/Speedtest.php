<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class Speedtest extends Model
{
    use HasFactory;

    /**
     * The attributes that are mass assignable.
     *
     * @var array
     */
    protected $fillable = [
        'ping',
        'download',
        'upload',
        'created_at',
        'server_id',
        'server_name',
        'server_host',
        'url',
        'scheduled',
        'failed',
    ];

    protected $table = 'speedtests';

    public function formatForInfluxDB()
    {
        return [
            'id' => (int) $this->id,
            'download' => (float) $this->download,
            'upload' => (float) $this->upload,
            'ping' => (float) $this->ping,
            'server_id' => (int) $this->server_id,
            'server_host' => $this->server_host,
            'server_name' => $this->server_name,
            'scheduled' => (bool) $this->scheduled,
        ];
    }
}
