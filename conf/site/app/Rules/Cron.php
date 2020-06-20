<?php

namespace App\Rules;

use Cron\CronExpression;
use Illuminate\Contracts\Validation\Rule;

class Cron implements Rule
{
    /**
     * Create a new rule instance.
     *
     * @return void
     */
    public function __construct()
    {
        //
    }

    /**
     * Determine if the value is a valid CRON expression
     *
     * @param  string  $attribute
     * @param  mixed  $value
     * @return bool
     */
    public function passes($attribute, $value)
    {
        return CronExpression::isValidExpression($value);
    }

    /**
     * Get the validation error message.
     *
     * @return string
     */
    public function message()
    {
        return 'The :attribute field must be a valid cron expression.';
    }
}
