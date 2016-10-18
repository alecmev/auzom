import callAPI from './api';

export function createOTP(body, onSuccess) { // eslint-disable-line
  return callAPI({
    url: '/otps',
    method: 'POST',
    body,
    type: 'OTP__CREATE',
    onSuccess,
  });
}
