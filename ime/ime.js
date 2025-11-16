importClass(android.os.Binder);
importClass(android.content.ClipData);
importClass(android.os.ServiceManager);
importClass(android.content.IClipboard);
importClass(android.util.Base64);

var FakeContext = {
    PACKAGE_NAME: "com.android.shell",
    ROOT_UID: 0
};

// 剪切板服务
var clipboardService = null;

// 获取剪切板服务
function getClipboardService() {
    if (clipboardService == null) {
        try {
            var binder = ServiceManager.getService("clipboard");
            if (binder == null) {
                return null;
            }

            clipboardService = android.content.IClipboard.Stub.asInterface(binder);
            if (clipboardService == null) {
                return null;
            }

        } catch (e) {
            return null;
        }
    }
    return clipboardService;
}

// 获取剪切板文本
function getClipboardText() {
    try {
        var service = getClipboardService();
        if (service == null) {
            return null;
        }

        var clipData = null;

        try {
            clipData = service.getPrimaryClip(FakeContext.PACKAGE_NAME);
        } catch (e1) {
            try {
                clipData = service.getPrimaryClip(FakeContext.PACKAGE_NAME, FakeContext.ROOT_UID);
            } catch (e2) {
                try {
                    clipData = service.getPrimaryClip(FakeContext.PACKAGE_NAME, null, FakeContext.ROOT_UID);
                } catch (e3) {
                    try {
                        clipData = service.getPrimaryClip(FakeContext.PACKAGE_NAME, null, FakeContext.ROOT_UID, 0);
                    } catch (e4) {
                        try {
                            clipData = service.getPrimaryClip(FakeContext.PACKAGE_NAME, FakeContext.ROOT_UID, null);
                        } catch (e5) {
                            try {
                                clipData = service.getPrimaryClip(FakeContext.PACKAGE_NAME, null, FakeContext.ROOT_UID, 0, true);
                            } catch (e6) {
                                return null;
                            }
                        }
                    }
                }
            }
        }

        if (clipData == null || clipData.getItemCount() == 0) {
            return null;
        }

        var item = clipData.getItemAt(0);
        var text = item.getText();
        return text != null ? text.toString() : null;
    } catch (e) {

    }
    return null;
}

// 设置剪切板文本
function setClipboardText(b64text) {
    if (b64text == null) {
        return false;
    }

    var decodedText;
    try {
        var decodedBytes = Base64.decode(b64text, Base64.DEFAULT);
        var javaString = new java.lang.String(decodedBytes, "UTF-8");
        decodedText = String(javaString);
    } catch (e) {
        return false;
    }

    try {
        var service = getClipboardService();
        if (service == null) {
            return false;
        }

        var clipData = ClipData.newPlainText(null, decodedText);

        try {
            service.setPrimaryClip(clipData, FakeContext.PACKAGE_NAME);
            return true;
        } catch (e1) {
            try {
                service.setPrimaryClip(clipData, FakeContext.PACKAGE_NAME, FakeContext.ROOT_UID);
                return true;
            } catch (e2) {
                try {
                    service.setPrimaryClip(clipData, FakeContext.PACKAGE_NAME, null, FakeContext.ROOT_UID);
                    return true;
                } catch (e3) {
                    try {
                        service.setPrimaryClip(clipData, FakeContext.PACKAGE_NAME, null, FakeContext.ROOT_UID, 0);
                        return true;
                    } catch (e4) {
                        try {
                            service.setPrimaryClip(clipData, FakeContext.PACKAGE_NAME, null, FakeContext.ROOT_UID, 0, true);
                            return true;
                        } catch (e5) {
                            return false;
                        }
                    }
                }
            }
        }
    } catch (e) {

    }
    return false;
}