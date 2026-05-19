Act as a Senior Golang Engineer. Please write unit tests for the Go service file I will provide at the end of this prompt.

Use the following strict standards and rules to ensure consistency with my existing codebase:

1. **Framework & Libraries:** Use the built-in `testing` package, `github.com/stretchr/testify/mock` for mocking, and `github.com/stretchr/testify/assert` for assertions.
2. **Design Pattern:** Implement Table-Driven Tests (using an anonymous slice of structs to iterate through test cases).
3. **Mock Handling:**
   - I ALREADY HAVE a separate `mocks_test.go` file containing the mock struct definitions (e.g., `MockUserRepo`, `MockActivityRepo`, `MockCache`).
   - DO NOT recreate these mock struct definitions in the generated test file. Initialize them directly using `new(MockRepoName)`.
   - Use mock assertion features like `.On()`, `.Return()`, `.Run()`, and `.Once()`.
   - Reset mock expectations inside the `t.Run` loop using `mockRepo.ExpectedCalls = nil` to prevent state leakage between test cases.
   - Call `mockRepo.AssertExpectations(t)` at the end of each test case.
4. **Casbin Handling:** If the service uses Casbin, initialize the enforcer in the setup function using `casbinModel.NewModelFromString(...)` with a standard RBAC model, and use the `&DummyCasbinAdapter{}` which I have already defined in my `mocks_test.go`. DO NOT leave the adapter as `nil`.
5. **Environment Variables & Utils:**
   - If the service calls encryption/hashing utilities (like `utils.Encrypt` or `utils.HashIndex`), set dummy environment variables in the test setup function using `os.Setenv("AES_256_KEY", "<base64_encoded_32_bytes_string>")` and `os.Setenv("PEPPER", "dummy-pepper")`.
   - You MUST use `defer os.Unsetenv(...)` at the beginning of each main test function to clean up the environment variables.
6. **Logging:** Discard logger output using `slog.New(slog.NewTextHandler(io.Discard, nil))` so it doesn't clutter the terminal during test execution.
7. **Test Cases:** Create both success scenarios (happy path) and failure scenarios (errors/edge cases) for each function being tested.

Here is my service code:

[PASTE_YOUR_SERVICE_CODE_HERE]

Please write the unit tests for the following functions: [SPECIFY_THE_FUNCTIONS_TO_TEST, e.g., GetAllUser, UpdateUser, DeleteUser]
