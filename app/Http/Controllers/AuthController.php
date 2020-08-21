<?php

namespace App\Http\Controllers;

use App\Auth\EmailVerification;
use App\Auth\LoginSession as AuthLoginSession;
use App\Helpers\EmailVerificationHelper;
use Illuminate\Support\Facades\Auth;
use App\Http\Controllers\Controller;
use App\LoginSession;
use App\Rules\CurrentPasswordMatches;
use App\User;
use DateTime;
use Hash;
use Illuminate\Support\Facades\Request as RequestFacade;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Validator;
use Log;
use Ramsey\Uuid\Uuid;

class AuthController extends Controller
{
    /**
     * Create a new AuthController instance.
     *
     * @return void
     */
    public function __construct()
    {
        $this->middleware('auth:api', ['except' => ['login', 'register']]);
    }

    public function register(Request $request)
    {
        $validator = Validator::make($request->all(), User::registerRules());

        if($validator->fails()) {
            return response()->json([
                'error' => $validator->errors()
            ], 422);
        }

        $user = User::create([
            'id' => Uuid::uuid4(),
            'name' => $request->name,
            'email' => $request->email,
            'password' => $request->password,
        ]);

        EmailVerification::create([
            'user_id' => $user->id,
            'token' => UUid::uuid4(),
            'expires' => new DateTime('+ 1 day')
        ]);

        $token = auth()->login($user);

        return $this->respondWithToken($token);
    }

    /**
     * Get a JWT via given credentials.
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function login(Request $request)
    {
        $validator = Validator::make($request->all(), User::loginRules());

        if($validator->fails()) {
            return response()->json([
                'error' => $validator->errors()
            ], 422);
        }

        $credentials = request(['email', 'password']);

        $length = 1440 * env('REMEMBER_TOKEN', 30);
        if (! $token = auth()->setTTL($length)->attempt($credentials)) {
            return response()->json(['error' => 'Unauthorized'], 401);
        }

        return $this->respondWithToken($token, $length);
    }

    /**
     * Get the authenticated User.
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function me()
    {
        return response()->json(auth()->user());
    }

    /**
     * Log the user out (Invalidate the token).
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function logout(Request $request)
    {
        $session = auth()->user()->sessions()->where('token', $request->token)->first();
        $session->active = false;
        $session->save();

        auth()->logout();

        return response()->json(['message' => 'Successfully logged out']);
    }

    /**
     * Refresh a token.
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function refresh(Request $request)
    {
        $session = auth()->user()->sessions()->where('token', $request->token)->first();
        $session->active = false;
        $session->save();

        return $this->respondWithToken(auth()->refresh());
    }

    /**
     * Get the token array structure.
     *
     * @param  string $token
     *
     * @return \Illuminate\Http\JsonResponse
     */
    protected function respondWithToken($token, $expiresIn = false)
    {
        if(!$expiresIn) {
            $expiresIn = 1440 * env('REMEMBER_TOKEN', 30);
        }
        $expiresTime = date("Y-m-d H:i:s", time() + ($expiresIn * 60));

        $expiresUnix = time() + ($expiresIn * 60);

        $this->storeSession($token, $expiresUnix);

        return response()->json([
            'access_token' => $token,
            'expires_in' => $expiresIn,
            'expires_on' => $expiresUnix,
            'expires_on_readable' => $expiresTime
        ]);
    }

    function storeSession($token, $expires)
    {
        AuthLoginSession::create([
            'id' => Uuid::uuid4(),
            'token' => $token,
            'user_id' => auth()->user()->id,
            'expires' => $expires,
            'ip' => RequestFacade::ip()
        ]);
    }

    public function getSessions()
    {
        $sessions = auth()->user()->sessions()->where([
            [ 'active', true ],
            [ 'expires', '>', time() ]
        ])->get();

        $sessions = $sessions->map(function ($item) {
            return collect($item)->forget(['token']);
        });

        return response()->json([
            'method' => 'get auth sessions',
            'response' => $sessions
        ], 200);
    }

    public function deleteSession($id)
    {
        $session = AuthLoginSession::where('id', $id)->firstOrFail();
        $session->delete();

        return response()->json([
            'method' => 'delete a login sesison',
            'response' => $session->id
        ], 200);
    }

    public function verifyEmail(Request $request)
    {
        $rules = [
            'user_id' => [ 'string', 'required' ],
            'token' => [ 'string', 'required' ],
        ];

        $validator = Validator::make($request->all(), $rules);
        if($validator->fails()) {
            return response()->json([
                'error' => $validator->errors()
            ], 422);
        }

        $user = EmailVerificationHelper::checkVerificationAttempt($request->user_id, $request->token);

        if(!$user) {
            return response()->json([
                'error' => 'token nout found'
            ], 404);
        }

        $user->email_verified_at = new DateTime();
        $user->save();

        return response()->json([
            'method' => 'verify email address',
            'success' => true,
        ], 200);
    }

    public function changePassword(Request $request)
    {
        $rules = [
            'currentPassword' => [ 'string', 'required', new CurrentPasswordMatches() ],
            'newPassword' => [ 'required', 'string', 'confirmed', 'min:8' ],
            'logoutDevices' => [ 'required', 'bool' ]
        ];

        $validator = Validator::make($request->all(), $rules);
        if($validator->fails()) {
            return response()->json([
                'method' => 'reset password',
                'success' => false,
                'error' => $validator->errors()
            ], 403);
        }

        $user = Auth::user();

        $user->password = $request->newPassword;
        $user->save();

        if($request->logoutDevices == true) {
            AuthLoginSession::where('user_id', $user->id)->update([ 'active' => false ]);
        }

        return response()->json([
            'method' => 'reset password',
            'success' => true
        ], 200);
    }
}
