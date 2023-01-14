import {persons} from "../utils";
import {getOrCreatePool, performScheduledDeletions, scheduleDeletion, sql} from "../db";
import {createConversation} from "./conversation/create";
import {revertSaveDraft, saveDraft, sendSaveDraftMutation} from "./draft/save";
import {sendRemoveDraftMutation} from "./draft/remove";
import {createMessage} from "./message/create";

const clearDb = async () => {
    const pool = await getOrCreatePool();

    await pool.connect(async (connection) => {
        const personIds = Object.values(persons)

        await connection.query(
            sql.typeAlias('void')`
                DELETE
                FROM messenger.draft
                WHERE author_id IN (${sql.join(personIds, sql.fragment`, `)})
            `
        )
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

test("save draft", async () => {
    const conversationId = await createConversation(persons.alice, persons.bob);
    const messageId = await createMessage(persons.alice, conversationId, 'Hi')

    const res = await sendSaveDraftMutation(persons.alice, conversationId, messageId, 'Hi!');
    scheduleDeletion(() => revertSaveDraft(persons.alice, conversationId))

    expect(res.error).toBeFalsy();
    expect(res.body).toMatchSnapshot({
        data: {
            saveDraft: {
                conversationId: expect.any(String),
                messageId: expect.any(String),
            },
        },
    });
    expect(res.body).toMatchObject({
        data: {
            saveDraft: {
                conversationId,
                messageId,
            },
        },
    })
});

test("remove draft", async () => {
    const conversationId = await createConversation(persons.alice, persons.bob);
    await saveDraft(persons.alice, conversationId, null, 'Hi!');

    const res = await sendRemoveDraftMutation(persons.alice, conversationId);

    expect(res.error).toBeFalsy();
    expect(res.body).toMatchSnapshot();
});
