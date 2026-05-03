export const buildUrlWithParams = (url: string, params: Record<string, string | number>) => {
    const keys = Object.keys(params);
    if (!keys.length) {
        return url;
    }

    let formattedUrl = url;

    keys.forEach((key: string, index: number) => {
        formattedUrl += `${index == 0 ? "?" : "&"}${key}=${params[key]}`
    });

    return formattedUrl;
}