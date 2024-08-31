let localIdentifier;

function processFile(evt) {
  const fileInput = evt.target;

  if (localIdentifier) {
    sessionStorage.removeItem(localIdentifier);
  }

  localIdentifier = Libra.newIdentifier();

  const file = fileInput.files[0],
    reader = new FileReader();

  reader.onload = function (e) {
    const data = e.target.result,
      fileInfo = {
        name: file.name,
        type: file.type,
      },
      key = Libra.createKey();

    console.log("File info", fileInfo);

    const encoder = new TextEncoder();

    return Promise.all([
      Libra.encrypt(data, key),
      Libra.encrypt(encoder.encode(JSON.stringify(fileInfo)), key),
    ])
      .then(([encryptedData, encryptedFileInfo]) => {
        console.log(
          "The encrypted data is " +
            encryptedData.cipher.byteLength +
            " bytes long"
        );

        const encryptedFile = new File([encryptedData.cipher], "file", {
          type: "application/octet-stream",
        });

        const encryptedInfo = new File([encryptedFileInfo.cipher], "info", {
          type: "application/octet-stream",
        });

        const dataTransfer = new DataTransfer();

        dataTransfer.items.add(encryptedFile);
        dataTransfer.items.add(encryptedInfo);

        fileInput.files = dataTransfer.files;

        document.querySelector('form input[name="local-identifier"]').value =
          localIdentifier;

        sessionStorage.setItem(
          localIdentifier,
          Libra.serialize(key, encryptedData.iv, encryptedFileInfo.iv)
        );

        document.querySelector("form").submit();
      })
      .catch(console.error);
  };

  reader.readAsArrayBuffer(file);
}
