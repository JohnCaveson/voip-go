const JOIN_AUDIO_PHRASES = [
  'Start Yapping',
  'Hop In',
  'Join the Huddle',
  'Sound On',
  'Tune In',
  'Join the Yap Session',
  'Let\'s Vibes',
  'Enter the Voice Dimension',
  'Unmute Your Soul',
  'Join the Noise',
  'Activate Voice Mode',
  'Ready to Rumble',
  'Step Into the Sound Booth',
  'Plugged In',
  'Commence Yapping',
]

export function getRandomJoinPhrase(): string {
  return JOIN_AUDIO_PHRASES[Math.floor(Math.random() * JOIN_AUDIO_PHRASES.length)]
}

export { JOIN_AUDIO_PHRASES }
