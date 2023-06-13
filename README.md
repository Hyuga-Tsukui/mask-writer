# Mask Writer is Sensitive value writer for Masking sensitive value in log

Mask Writer is a confidential value writer to mask confidential values in logs, etc. It implements io.Writer and can be passed to anything that uses io.

However, it can only be used for structured logs as it is designed to handle JSON.