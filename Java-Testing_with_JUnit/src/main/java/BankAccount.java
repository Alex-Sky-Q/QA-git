public class BankAccount {
    private double balance;
    private double minBalance;
    private boolean isActive;
    private String holderName;

    public BankAccount(double balance, double minBalance) {
        this.balance = balance;
        this.minBalance = minBalance;
        isActive = true;
    }

    public boolean isActive() {
        return isActive;
    }

    public String getHolderName() {
        return holderName;
    }

    public void setHolderName(String holderName) {
        this.holderName = holderName;
    }

    public double getBalance() {
        return balance;
    }

    public double getMinBalance() {
        return minBalance;
    }

    public double deposit(double amount) {
        return balance += amount;
    }

    public double withdraw(double amount) {
        try {
            Thread.sleep(1);
        } catch (InterruptedException e) {
            e.printStackTrace();
        }
        if (balance - amount >= minBalance) {
            return balance -= amount;
        }
        else {
            throw new RuntimeException();
        }
    }

    public double splitAccountBalance(int parts) {
        return balance /= parts;
    }

}
