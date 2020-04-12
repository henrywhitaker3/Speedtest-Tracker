<?php

namespace App\Helpers;

use App\Auth\EmailVerification;
use App\User;

class EmailVerificationHelper {
    public static function checkVerificationAttempt($userID, $token)
    {
        $verification = EmailVerification::where('user_id', $userID)
                                         ->where('token', $token)
                                         ->first();

        if(!$verification) {
            return false;
        }

        return User::where('id', $userID)->first();
    }

    public static function verifyUser(User $user)
    {

    }

    public static function userIsVerified()
    {
        if(auth()->user()->email_verified_at == null) {
            return false;
        }

        return true;
    }
}
