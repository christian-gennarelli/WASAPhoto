export const getImgUrl = (path) => {
    if (path) {
        return __API_URL__ + "/photos/?photo_path=" + path
    }
}