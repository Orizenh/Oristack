import { fetchAPI } from "$lib"

export async function load():Promise<any>{
    return await fetchAPI("/")
}