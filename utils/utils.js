importClass(android.os.ServiceManager);
importClass(android.hardware.display.IDisplayManager$Stub);

var dm = null;

function initDisplayManager() {
    try {
        var binder = ServiceManager.getService("display");
        if (binder == null) {
            return;
        }
        dm = IDisplayManager$Stub.asInterface(binder);
    } catch (e) {

    }
}

function wmSize() {
    try {
        if (dm == null) {
            initDisplayManager();
        }
        var displayInfo = dm.getDisplayInfo(0);
        return displayInfo.logicalWidth + "x" + displayInfo.logicalHeight;
    } catch (e) {
        return e;
    }
}

wmSize();