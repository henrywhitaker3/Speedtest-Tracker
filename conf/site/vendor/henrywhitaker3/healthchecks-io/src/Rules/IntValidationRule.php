<?php

namespace Henrywhitaker3\Healthchecks\Rules;

use Rakit\Validation\Rule;

class IntValidationRule extends Rule
{
    protected $message = ':attribute :value has been used';

    public function __construct()
    {
        //
    }

    /**
     * Check the value is an int
     *
     * @param mixed $value
     * @return bool
     */
    public function check($value): bool
    {
        if(is_int($value)) {
            return true;
        }

        return false;
    }
}