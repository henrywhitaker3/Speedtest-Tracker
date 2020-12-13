<?php

namespace NotificationChannels\Telegram;

use Illuminate\Support\Facades\View;
use JsonSerializable;
use NotificationChannels\Telegram\Traits\HasSharedLogic;
use Psr\Http\Message\StreamInterface;

/**
 * Class TelegramFile.
 */
class TelegramFile implements JsonSerializable
{
    use HasSharedLogic;

    /** @var string content type. */
    public $type = 'document';

    /**
     * @param string $content
     *
     * @return self
     */
    public static function create(string $content = ''): self
    {
        return new self($content);
    }

    /**
     * Message constructor.
     *
     * @param string $content
     */
    public function __construct(string $content = '')
    {
        $this->content($content);
        $this->payload['parse_mode'] = 'Markdown';
    }

    /**
     * Notification message (Supports Markdown).
     *
     * @param string $content
     *
     * @return $this
     */
    public function content(string $content): self
    {
        $this->payload['caption'] = $content;

        return $this;
    }

    /**
     * Add File to Message.
     *
     * Generic method to attach files of any type based on API.
     *
     * @param string|resource|StreamInterface $file
     * @param string                          $type
     * @param string|null                     $filename
     *
     * @return $this
     */
    public function file($file, string $type, string $filename = null): self
    {
        $this->type = $type;

        $isLocalFile = $this->isReadableFile($file);

        if ($filename !== null || $isLocalFile) {
            $this->payload['file'] = [
                'filename' => $filename,
                'name'     => $type,
                'contents' => $isLocalFile ? fopen($file, 'rb') : $file,
            ];
        } else {
            $this->payload[$type] = $file;
        }

        return $this;
    }

    /**
     * Attach an image.
     *
     * Use this method to send photos.
     *
     * @param string $file
     *
     * @return $this
     */
    public function photo(string $file): self
    {
        return $this->file($file, 'photo');
    }

    /**
     * Attach an audio file.
     *
     * Use this method to send audio files, if you want Telegram clients to display them in the music player.
     * Your audio must be in the .mp3 format.
     *
     * @param string $file
     *
     * @return $this
     */
    public function audio(string $file): self
    {
        return $this->file($file, 'audio');
    }

    /**
     * Attach a document or any file as document.
     *
     * Use this method to send general files.
     *
     * @param string      $file
     * @param string|null $filename
     *
     * @return $this
     */
    public function document(string $file, string $filename = null): self
    {
        return $this->file($file, 'document', $filename);
    }

    /**
     * Attach a video file.
     *
     * Use this method to send video files, Telegram clients support mp4 videos.
     *
     * @param string $file
     *
     * @return $this
     */
    public function video(string $file): self
    {
        return $this->file($file, 'video');
    }

    /**
     * Attach an animation file.
     *
     * Use this method to send animation files (GIF or H.264/MPEG-4 AVC video without sound).
     *
     * @param string $file
     *
     * @return $this
     */
    public function animation(string $file): self
    {
        return $this->file($file, 'animation');
    }

    /**
     * Attach a voice file.
     *
     * Use this method to send audio files, if you want Telegram clients to display the file as a playable voice
     * message. For this to work, your audio must be in an .ogg file encoded with OPUS.
     *
     * @param string $file
     *
     * @return $this
     */
    public function voice(string $file): self
    {
        return $this->file($file, 'voice');
    }

    /**
     * Attach a video note file.
     *
     * Telegram clients support rounded square mp4 videos of up to 1 minute long.
     * Use this method to send video messages.
     *
     * @param string $file
     *
     * @return $this
     */
    public function videoNote(string $file): self
    {
        return $this->file($file, 'video_note');
    }

    /**
     * Attach a view file as the content for the notification.
     * Supports Laravel blade template.
     *
     * @param string $view
     * @param array  $data
     * @param array  $mergeData
     *
     * @return $this
     */
    public function view(string $view, array $data = [], array $mergeData = []): self
    {
        return $this->content(View::make($view, $data, $mergeData)->render());
    }

    /**
     * Determine there is a file.
     *
     * @return bool
     */
    public function hasFile(): bool
    {
        return isset($this->payload['file']);
    }

    /**
     * Returns params payload.
     *
     * @return array
     */
    public function toArray(): array
    {
        return $this->hasFile() ? $this->toMultipart() : $this->payload;
    }

    /**
     * Create Multipart array.
     *
     * @return array
     */
    public function toMultipart(): array
    {
        $data = [];
        foreach ($this->payload as $name => $contents) {
            $data[] = ($name === 'file') ? $contents : compact('name', 'contents');
        }

        return $data;
    }

    /**
     * Determine if it's a regular and readable file.
     *
     * @param string $file
     *
     * @return bool
     */
    protected function isReadableFile(string $file): bool
    {
        return is_file($file) && is_readable($file);
    }
}
