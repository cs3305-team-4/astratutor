export interface Detail {
    msg: string;
}

export interface RESTError {
    detail: Detail;
}

export interface RESTInfo {
    detail: Detail;
}

/**
 * Perform a fetch request, throw an exception with a pretty message if the fetch() function throws an exception or does not return the expected status
 * @param input Standard request info passed to fetch
 * @param init
 * @param expectedStatuses Throw an exception if the returned status code is not in this array
 */
export async function fetchRest(
    input: RequestInfo,
    init?: RequestInit | undefined,
    expectedStatuses: number[] = [200, 201, 202, 203, 204, 205, 206],
): Promise<Response> {
    try {
        const res = await fetch(input, init);

        if (!expectedStatuses.includes(res.status)) {
            let msg = `${res.status}: ${res.statusText}`;

            // Test if the error returned has a rest error payload
            try {
                const err: RESTError = await res.json();

                if (err?.detail?.msg) {
                    msg = err.detail.msg;
                } else {
                    msg = `${res.status}: ${res.statusText}`;
                }
            } catch {}

            // Throw message
            throw new Error(msg);
        }

        return res;
    } catch (e) {
        throw new Error(`${e}`);
    }
}
