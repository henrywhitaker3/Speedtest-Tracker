<?php

namespace Database\Factories;

use App\Models\Speedtest;
use Illuminate\Database\Eloquent\Factories\Factory;

class SpeedtestFactory extends Factory
{
    /**
     * The name of the factory's corresponding model.
     *
     * @var string
     */
    protected $model = Speedtest::class;

    /**
     * Define the model's default state.
     *
     * @return array
     */
    public function definition()
    {
        return [
            'download' => rand(15, 900),
            'upload' => rand(15, 900),
            'ping' => rand(1, 25),
            'scheduled' => (bool) rand(0, 1),
            'failed' => false,
        ];
    }
}
