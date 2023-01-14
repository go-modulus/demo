import {persons} from "../utils";
import {getOrCreatePool, performScheduledDeletions, scheduleDeletion, sql} from "../db";
import {NotFoundError} from "slonik";
import {
    createConversation,
    revertCreateConversation,
    sendCreateOneToOneConversationMutation
} from "./conversation/create";
import {sendMyConversationsQuery} from "./conversation/my";

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

test("create conversation", async () => {
    const res = await sendCreateOneToOneConversationMutation(persons.alice, persons.iggy);

    expect(res.error).toBeFalsy();
    expect(res.body).toMatchSnapshot({
        data: {
            createOneToOneConversation: {
                id: expect.any(String),
                createdAt: expect.any(Number),
            }
        }
    });

    const conversationId = res.body.data?.createOneToOneConversation.id;
    scheduleDeletion(() => revertCreateConversation(conversationId))
});

test("list conversations", async () => {
    const conversationIds = await Promise.all([
        createConversation(persons.alice, persons.bob),
        createConversation(persons.alice, persons.dawid),
        createConversation(persons.iggy, persons.alice),
    ])

    const res = await sendMyConversationsQuery(persons.alice);

    expect(res.error).toBeFalsy();
    expect(res.body).toMatchSnapshot({
        data: {
            myConversations: {
                edges: conversationIds.map(() => ({
                    cursor: expect.any(String),
                    node: {
                        id: expect.any(String),
                        createdAt: expect.any(Number),
                    },
                })),
                pageInfo: {
                    startCursor: expect.any(String),
                    endCursor: expect.any(String),
                },
            },
        },
    });

    const conversationIdsInResponse = res.body.data.myConversations.edges.map(
        (edge: any) => edge.node.id
    );

    expect(conversationIdsInResponse).toIncludeAllMembers(conversationIds);
});