<?php

namespace App\Casts;

use Illuminate\Contracts\Database\Eloquent\CastsAttributes;

class CommaSeparatedArrayCast implements CastsAttributes
{
    /**
     * Array of settings that should be cast
     */
    private array $shouldCast = [
        'visible_columns',
        'hidden_columns',
    ];

    /**
     * Cast the given value.
     *
     * @param  \Illuminate\Database\Eloquent\Model  $model
     * @param  string  $key
     * @param  mixed  $value
     * @param  array  $attributes
     * @return mixed
     */
    public function get($model, $key, $value, $attributes)
    {
        if (!in_array($model->name, $this->shouldCast)) {
            return $value;
        }

        return explode(',', $value);
    }

    /**
     * Prepare the given value for storage.
     *
     * @param  \Illuminate\Database\Eloquent\Model  $model
     * @param  string  $key
     * @param  mixed  $value
     * @param  array  $attributes
     * @return mixed
     */
    public function set($model, $key, $value, $attributes)
    {
        if (!in_array($model->name, $this->shouldCast)) {
            return $value;
        }

        return implode(',', $value);
    }
}
