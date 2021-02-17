/**
 * Get array of connected media devices
 * @returns {Promise<MediaDeviceInfo[]>} media devices
 */
export async function getDevices(): Promise<MediaDeviceInfo[]> {
  return await window.navigator.mediaDevices.enumerateDevices();
}

/**
 * Check microphone permission
 * @returns {Promise<boolean>} true if microphone acces granted
 */
export async function microphonePermission(): Promise<boolean> {
  const devices = await getDevices();
  return (
    devices.filter((device) => {
      return device.kind === 'audioinput' && device.label !== 'default' && device.label !== '';
    }).length > 0
  );
}

/**
 * Check camera permission
 * @returns {Promise<boolean>} true if camera acces granted
 */
export async function cameraPermission(): Promise<boolean> {
  const devices = await getDevices();
  return (
    devices.filter((device) => {
      return device.kind === 'videoinput' && device.label !== 'default' && device.label !== '';
    }).length > 0
  );
}

/**
 * Prompt User for access to Mic + Camera
 */
export async function devicePermissions() {
  try {
    console.log('Attempting to get Camera and Microphone Permission');
    await navigator.mediaDevices.getUserMedia({ audio: true, video: true });
  } catch (err) {
    console.error(err);
    try {
      console.log('Attempting to get Microphone Permission');
      await navigator.mediaDevices.getUserMedia({ audio: true });
    } catch (err) {
      console.error(err);
      try {
        console.log('Attempting to get Camera Permission');
        await navigator.mediaDevices.getUserMedia({ video: true });
      } catch (err) {
        console.error(err);
      }
    }
  }
}

/**
 * Get the Audio Stream from *device*
 * @param {string} deviceId ID of device to access
 * @returns {Promise<(MediaStream | null)>} audio stream
 */
export async function audioStream(deviceId: string): Promise<MediaStream | null> {
  try {
    return await navigator.mediaDevices.getUserMedia({ audio: { deviceId: deviceId } });
  } catch (err) {
    console.error(err);
  }
  return null;
}

/**
 * Get the Camera Stream from *device*
 * @param {string} deviceId ID of device to access
 * @returns {Promise<(MediaStream | null)>} camera stream
 */
export async function cameraStream(deviceId: string, audioId?: string): Promise<MediaStream | null> {
  try {
    return await navigator.mediaDevices.getUserMedia({ video: { deviceId: deviceId }, audio: { deviceId: audioId } });
  } catch (err) {
    console.error(err);
  }
  return null;
}

export async function screenStream(): Promise<MediaStream | null> {
  try {
    // https://github.com/microsoft/TypeScript/issues/33232#issuecomment-616131983
    const mediaDevices = navigator.mediaDevices as any;
    return await mediaDevices.getDisplayMedia({ audio: true, video: true });
  } catch (err) {
    console.error(err);
  }
  return null;
}
