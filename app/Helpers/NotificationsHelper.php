<?php

namespace App\Helpers;

class NotificationsHelper {

    /**
     * Parse $errors and format message
     *
     * @param array $errors
     * @return String
     */
    public static function formatPercentageThresholdMessage(array $errors)
    {
        $msg = NotificationsHelper::thresholdMessageStart($errors);

        $msg = $msg . 'exceeded the percentage threshold';

        return $msg;
    }

    /**
     * Parse $errors and format message
     *
     * @param array $errors
     * @return String
     */
    public static function formatAbsoluteThresholdMessage(array $errors)
    {
        $msg = NotificationsHelper::thresholdMessageStart($errors);

        $msg = $msg . 'exceeded the absolute threshold';

        return $msg;
    }

    /**
     * Iterate through errors to format message
     *
     * @param array $errors
     * @return String
     */
    public static function thresholdMessageStart(array $errors)
    {
        $msg = 'For the latest speedtest, the ';

        for($i = 0; $i < sizeof($errors); $i++) {
            $key = $errors[$i];
            $msg = $msg . $key;
            if(sizeof($errors) > 1 && $i < (sizeof($errors) - 1)) {
                $msg = $msg . ', ';
            }
        }

        if($msg[-1] != '') {
            $msg = $msg . ' ';
        }

        if(sizeof($errors) > 1) {
            $msg = $msg . 'values ';
        } else {
            $msg = $msg . 'value ';
        }

        return $msg;
    }
}
