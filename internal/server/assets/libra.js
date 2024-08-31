(function (Libra) {
  const ALGORITHM = "AES-GCM";
  const SIZE = 256;
  const KEY_USAGES = ["encrypt", "decrypt"];

  Libra.newIdentifier = function () {
    return crypto.randomUUID();
  };

  Libra.createKey = function () {
    return crypto.getRandomValues(new Uint8Array(SIZE / 8));
  };

  Libra.serialize = function (key, ...ivs) {
    const str =
      btoa(String.fromCharCode(...key)) +
      "." +
      ivs.map((iv) => btoa(String.fromCharCode(...iv))).join(".");
    return str;
  };

  Libra.deserialize = function (raw) {
    const parts = raw.split(".");
    return {
      key: base64ToArrayBuffer(parts[0].replace(/^#/, "")),
      ivs: parts.slice(1).map((iv) => base64ToArrayBuffer(iv)),
    };
  };

  function base64ToArrayBuffer(base64) {
    const binaryString = atob(base64);
    const bytes = new Uint8Array(binaryString.length);
    for (let i = 0; i < binaryString.length; i++) {
      bytes[i] = binaryString.charCodeAt(i);
    }
    return bytes;
  }

  Libra.encrypt = function (data, key) {
    const iv = crypto.getRandomValues(new Uint8Array(12));
    return crypto.subtle
      .importKey(
        "raw",
        key,
        {
          name: ALGORITHM,
          length: SIZE,
        },
        true,
        KEY_USAGES
      )
      .then((key) => {
        return crypto.subtle.encrypt(
          {
            name: ALGORITHM,
            iv,
          },
          key,
          data
        );
      })
      .then((cipher) => ({ cipher, iv }));
  };

  Libra.decrypt = function (cipher, key, iv) {
    return crypto.subtle
      .importKey(
        "raw",
        key,
        {
          name: ALGORITHM,
          length: SIZE,
        },
        true,
        KEY_USAGES
      )
      .then((key) => {
        console.log(cipher, key, iv);
        return crypto.subtle.decrypt(
          {
            name: ALGORITHM,
            iv,
          },
          key,
          cipher
        );
      });
  };
})((Libra = {}));
