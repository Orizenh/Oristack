import { env } from "$env/dynamic/private"


export async function fetchAPI(
    path: string, 
    method?: string|undefined, 
    data?: Array<any>|undefined
):Promise<any>{
    if(method === undefined){
        method = "GET"
    }
    let headers:RequestInit = {
        method: method
    }
    if(data !== undefined){
        headers.body = JSON.stringify(data)
    }
    let result
    let response
    try{
        result = await fetch(env.API_URL+path, headers)
        response = await result.json()
        if(!result.ok){
            throw response.error
        }
    }catch(error:any){
        return { error: error }
    }
    return {
        data: response
    }
}