<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

class CreateActiveSessionsTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::create('active_sessions', function (Blueprint $table) {
            $table->uuid('id')->primary();
            $table->text('token');
            $table->boolean('active')->default(true);
            $table->uuid('user_id');
            $table->bigInteger('expires');
            $table->string('ip');
            $table->timestamps();
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('active_sessions');
    }
}
