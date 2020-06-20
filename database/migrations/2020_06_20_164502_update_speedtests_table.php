<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

class UpdateSpeedtestsTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::table('speedtests', function($table) {
            $table->integer('server_id')->nullable();
            $table->string('server_name')->nullable();
            $table->string('server_host')->nullable();
            $table->string('url')->nullable();
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::table('speedtests', function($table) {
            $table->dropColumn('server_id');
        });

        Schema::table('speedtests', function($table) {
            $table->dropColumn('server_name');
        });

        Schema::table('speedtests', function($table) {
            $table->dropColumn('server_host');
        });

        Schema::table('speedtests', function($table) {
            $table->dropColumn('url');
        });
    }
}
