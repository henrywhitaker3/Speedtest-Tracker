<?php

namespace Tests\Feature;

use App\Rules\Cron;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Illuminate\Foundation\Testing\WithFaker;
use Tests\TestCase;

class CronRuleTest extends TestCase
{
    /**
     * Test a valid CRON expression
     *
     * @return void
     */
    public function testValidCronValidationRule()
    {
        $rule = [
            'test' => new Cron,
        ];

        $data = [
            'test' => '*/5 * * * *',
        ];

        $validator = $this->app['validator']->make($data, $rule);
        $this->assertTrue($validator->passes());
    }

    /**
     * Test an invalid CRON expression
     *
     * @return void
     */
    public function testInvalidCronValidationRule()
    {
        $rule = [
            'test' => new Cron,
        ];

        $data = [
            'test' => 'invalid',
        ];

        $validator = $this->app['validator']->make($data, $rule);
        $this->assertFalse($validator->passes());
    }
}
