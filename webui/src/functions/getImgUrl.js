export const getImgUrl = (path) => {
    if (path) {
        return "http://localhost:3000/photos/?photo_path=" + path
    }
}