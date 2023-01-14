import {persons} from "../utils";
import {getOrCreatePool, performScheduledDeletions, scheduleDeletion, sql} from "../db";
import {NotFoundError} from "slonik";
import {createMessage, revertCreateMessage, sendCreateMessageMutation} from "./message/create";
import {sendEditMessage} from "./message/edit";
import {createConversation} from "./conversation/create";

const clearDb = async () => {
    const pool = await getOrCreatePool();

    await pool.connect(async (connection) => {
        const personIds = Object.values(persons)

        try {
            const results = await connection.many(
                sql.typeAlias('id')`
                    SELECT id
                    FROM messenger.conversation
                    WHERE sender_id IN (${sql.join(personIds, sql.fragment`, `)})
                       OR receiver_id IN (${sql.join(personIds, sql.fragment`, `)})
                `
            )
            const conversationIds = results.map(({id}) => id)

            await connection.query(
                sql.typeAlias('void')`
                    DELETE
                    FROM messenger.message
                    WHERE conversation_id IN (${sql.join(conversationIds, sql.fragment`, `)})
                `
            )

            await connection.query(
                sql.typeAlias('void')`
                    DELETE
                    FROM messenger.conversation
                    WHERE id IN (${sql.join(conversationIds, sql.fragment`, `)})
                `
            )
        } catch (e) {
            if (!(e instanceof NotFoundError)) {
                throw e;
            }
        }
    })
}

beforeAll(clearDb)
afterEach(performScheduledDeletions)
afterAll(async () => {
    const pool = await getOrCreatePool();
    if (pool.getPoolState().ended) {
        return;
    }

    await pool.end();
})

test("create message", async () => {
    const conversationId = await createConversation(persons.alice, persons.bob);

    const res = await sendCreateMessageMutation(persons.alice, conversationId, 'Hi!');

    expect(res.error).toBeFalsy();
    expect(res.body).toMatchSnapshot({
        data: {
            createMessage: {
                id: expect.any(String),
                conversationId: expect.any(String),
                updatedAt: expect.any(Number),
                createdAt: expect.any(Number),
            }
        }
    });
    expect(res.body).toMatchObject({
        data: {
            createMessage: {
                conversationId,
            },
        },
    })

    const messageId = res.body.data.createMessage.id;
    scheduleDeletion(() => revertCreateMessage(messageId))
});

test("edit message", async () => {
    const conversationId = await createConversation(persons.dawid, persons.bob);
    const messageId = await createMessage(persons.dawid, conversationId, 'Hi!');

    const res = await sendEditMessage(persons.dawid, messageId, 'Hello!');

    expect(res.error).toBeFalsy();
    expect(res.body).toMatchSnapshot({
        data: {
            editMessage: {
                id: expect.any(String),
                conversationId: expect.any(String),
                updatedAt: expect.any(Number),
                createdAt: expect.any(Number),
            },
        },
    });
    expect(res.body).toMatchObject({
        data: {
            editMessage: {
                id: messageId,
                conversationId,
            },
        },
    })
});
