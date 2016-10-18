export function messagePush(text, isError = false, code, url) {
  return {
    type: 'MESSAGE_PUSH',
    payload: {
      text,
      isError,
      code,
      url,
    },
  };
}

export function messageShift() {
  return {
    type: 'MESSAGE_SHIFT',
  };
}
