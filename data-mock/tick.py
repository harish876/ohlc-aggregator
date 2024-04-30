import pandas as pd
import numpy as np

# Generate sample tick data between 9:01 and 9:02
start_time = pd.Timestamp('2024-04-18 09:01:00')
end_time = pd.Timestamp('2024-04-18 09:02:00')

# Generate timestamps between 9:01 and 9:02 with variability
num_ticks = 5000
timestamps = pd.date_range(start=start_time, end=end_time, periods=num_ticks)
print(timestamps)
ltt_seconds = np.random.randint(0, 60, size=num_ticks)
ltt = timestamps + pd.to_timedelta(ltt_seconds, unit='s')

# Generate random LTP values
ltp = np.random.randint(100, 200, size=num_ticks)

# Create DataFrame with LTP and LTT
tick_data = pd.DataFrame({'LTP': ltp}, index=ltt)
print(tick_data)

# Resample tick data to every minute and calculate OHLC using ohlc() method
ohlc_data = tick_data.resample('1min').ohlc()

print(ohlc_data)
