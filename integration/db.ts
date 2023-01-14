import {createPool as createSlonik, createSqlTag, DatabasePool} from "slonik";
import {z} from "zod";

export const sql = createSqlTag({
    typeAliases: {
        id: z.object({
            id: z.string(),
        }),
        void: z.object({}).strict(),
    }
})

let pool: DatabasePool | null;
let scheduledDeletions: Array<() => Promise<void>> = []

export const scheduleDeletion = (performDeletion: () => Promise<void>) => {
    scheduledDeletions.push(performDeletion)
}

export const performScheduledDeletions = async () => {
    do {
        const performDeletion = scheduledDeletions.pop()

        if (performDeletion === undefined) {
            break
        }

        await performDeletion();
    } while (scheduledDeletions.length > 0)
}

export const getOrCreatePool = async () => {
    if (pool) {
        return pool;
    }

    return pool = await createSlonik('postgres://modulus:secret@localhost:5432/demo');
}